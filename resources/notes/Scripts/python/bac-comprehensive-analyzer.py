#!/usr/bin/env python3
"""
BAC Comprehensive Analyzer
Analyzes all BAC exams to identify patterns, topics, and create study priorities
"""

import os
import re
from collections import Counter, defaultdict
from pathlib import Path
import json

class BACAnalyzer:
    def __init__(self, exams_dir):
        self.exams_dir = Path(exams_dir)
        self.topics = defaultdict(int)
        self.keywords = Counter()
        self.exercises_by_type = defaultdict(list)
        self.french_vocab = set()
        
    def analyze_all_exams(self):
        """Analyze all BAC exam files"""
        results = {
            'total_exams': 0,
            'topics_frequency': {},
            'exercise_patterns': {},
            'difficulty_analysis': {},
            'french_vocabulary': [],
            'study_priorities': []
        }
        
        exam_files = list(self.exams_dir.glob('BAC-*.md'))
        results['total_exams'] = len(exam_files)
        
        for exam_file in exam_files:
            self._analyze_exam(exam_file)
        
        # Compile results
        results['topics_frequency'] = dict(sorted(self.topics.items(), key=lambda x: x[1], reverse=True))
        results['exercise_patterns'] = self._identify_patterns()
        results['french_vocabulary'] = sorted(self.french_vocab)
        results['study_priorities'] = self._generate_priorities()
        
        return results
    
    def _analyze_exam(self, exam_path):
        """Analyze a single exam file"""
        with open(exam_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Extract year and session
        year_match = re.search(r'BAC-(\d{4})', exam_path.name)
        year = year_match.group(1) if year_match else 'unknown'
        
        # Identify topics
        self._extract_topics(content)
        
        # Extract French vocabulary
        self._extract_french_vocab(content)
        
        # Categorize exercises
        self._categorize_exercises(content, year)
    
    def _extract_topics(self, content):
        """Extract and count topics from exam content"""
        # Chemistry topics
        chem_topics = {
            'cinétique': r'cinétique|vitesse.*réaction|catalyseur|demi-vie',
            'redox': r'oxydation|réduction|redox|oxydant|réducteur',
            'acide-base': r'acide|base|pH|pKa|tampon|titration|indicateur',
            'organique': r'alcool|aldéhyde|cétone|ester|amine|acide.*carboxylique',
            'équilibre': r'équilibre.*chimique|constante.*équilibre|Le.*Chatelier'
        }
        
        # Physics topics
        phys_topics = {
            'mécanique': r'énergie.*mécanique|ressort|oscillation|projectile|friction',
            'électromagnétisme': r'champ.*électrique|champ.*magnétique|induction|Lorentz|flux',
            'circuits-AC': r'RLC|résonance|impédance|déphasage|Fresnel',
            'bobine': r'bobine|auto-induction|f\.e\.m|inductance',
            'ondes': r'onde|fréquence|longueur.*onde|interférence'
        }
        
        all_topics = {**chem_topics, **phys_topics}
        
        for topic, pattern in all_topics.items():
            matches = len(re.findall(pattern, content, re.IGNORECASE))
            if matches > 0:
                self.topics[topic] += matches
    
    def _extract_french_vocab(self, content):
        """Extract French technical vocabulary"""
        # Common technical terms in BAC exams
        patterns = [
            r'\b[A-ZÀÂÄÇÉÈÊËÏÎÔÙÛÜ][a-zàâäçéèêëïîôùûü]+(?:\s+[a-zàâäçéèêëïîôùûü]+)*\b',
        ]
        
        # Extract sentences with technical terms
        sentences = re.split(r'[.!?]', content)
        for sentence in sentences:
            # Look for French technical vocabulary
            words = re.findall(r'\b[a-zàâäçéèêëïîôùûüA-ZÀÂÄÇÉÈÊËÏÎÔÙÛÜ]{4,}\b', sentence)
            for word in words:
                if any(c in word.lower() for c in 'àâäçéèêëïîôùûü'):
                    self.french_vocab.add(word.lower())
    
    def _categorize_exercises(self, content, year):
        """Categorize exercises by type"""
        exercises = re.findall(r'Exercice\s+(\d+)(.*?)(?=Exercice\s+\d+|Corrigé|$)', content, re.DOTALL | re.IGNORECASE)
        
        for ex_num, ex_content in exercises:
            ex_type = self._identify_exercise_type(ex_content)
            self.exercises_by_type[ex_type].append({
                'year': year,
                'number': ex_num,
                'preview': ex_content[:200].strip()
            })
    
    def _identify_exercise_type(self, content):
        """Identify the type of exercise"""
        content_lower = content.lower()
        
        if any(word in content_lower for word in ['cinétique', 'vitesse', 'catalyseur']):
            return 'Chimie - Cinétique'
        elif any(word in content_lower for word in ['acide', 'base', 'ph', 'pka']):
            return 'Chimie - Acide-Base'
        elif any(word in content_lower for word in ['alcool', 'aldéhyde', 'organique']):
            return 'Chimie - Organique'
        elif any(word in content_lower for word in ['ressort', 'énergie', 'oscillation']):
            return 'Physique - Mécanique'
        elif any(word in content_lower for word in ['bobine', 'induction', 'champ']):
            return 'Physique - Électromagnétisme'
        elif any(word in content_lower for word in ['rlc', 'résonance', 'impédance']):
            return 'Physique - Circuits AC'
        else:
            return 'Autre'
    
    def _identify_patterns(self):
        """Identify recurring patterns in exercises"""
        patterns = {}
        for ex_type, exercises in self.exercises_by_type.items():
            patterns[ex_type] = {
                'count': len(exercises),
                'frequency': f"{len(exercises)/22*100:.1f}%",  # 22 total exams
                'appears_in': len(set(ex['year'] for ex in exercises))
            }
        return patterns
    
    def _generate_priorities(self):
        """Generate study priorities based on frequency"""
        priorities = []
        
        # Sort topics by frequency
        sorted_topics = sorted(self.topics.items(), key=lambda x: x[1], reverse=True)
        
        for i, (topic, count) in enumerate(sorted_topics[:10], 1):
            priority = {
                'rank': i,
                'topic': topic,
                'frequency': count,
                'priority': 'HIGH' if i <= 3 else 'MEDIUM' if i <= 6 else 'NORMAL',
                'recommendation': self._get_recommendation(topic, count)
            }
            priorities.append(priority)
        
        return priorities
    
    def _get_recommendation(self, topic, frequency):
        """Get study recommendation for a topic"""
        if frequency >= 20:
            return f"CRITICAL: Appears {frequency} times. Master this topic completely."
        elif frequency >= 10:
            return f"IMPORTANT: Appears {frequency} times. Focus heavily on this."
        else:
            return f"MODERATE: Appears {frequency} times. Review regularly."


def main():
    exams_dir = Path(__file__).parent / 'Study-Vault/04-Exams/BAC-2002-2012'
    
    if not exams_dir.exists():
        print(f"Error: Directory not found: {exams_dir}")
        return
    
    print("🔍 Analyzing BAC Exams (2002-2012)...")
    print("=" * 60)
    
    analyzer = BACAnalyzer(exams_dir)
    results = analyzer.analyze_all_exams()
    
    # Save results
    output_file = Path(__file__).parent / 'BAC-ANALYSIS-REPORT.json'
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(results, f, ensure_ascii=False, indent=2)
    
    # Print summary
    print(f"\n📊 ANALYSIS COMPLETE")
    print(f"Total Exams Analyzed: {results['total_exams']}")
    print(f"\n🎯 TOP 10 PRIORITY TOPICS:")
    print("-" * 60)
    
    for priority in results['study_priorities']:
        print(f"{priority['rank']}. [{priority['priority']}] {priority['topic']}")
        print(f"   Frequency: {priority['frequency']} occurrences")
        print(f"   {priority['recommendation']}")
        print()
    
    print(f"\n📝 EXERCISE PATTERNS:")
    print("-" * 60)
    for ex_type, data in sorted(results['exercise_patterns'].items(), key=lambda x: x[1]['count'], reverse=True):
        print(f"{ex_type}: {data['count']} exercises ({data['frequency']}) - Appears in {data['appears_in']} years")
    
    print(f"\n💾 Full report saved to: {output_file}")
    print(f"📚 French vocabulary extracted: {len(results['french_vocabulary'])} terms")

if __name__ == '__main__':
    main()
