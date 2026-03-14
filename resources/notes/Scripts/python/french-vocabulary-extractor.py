#!/usr/bin/env python3
"""
French Vocabulary Extractor for BAC Exams
Extracts and translates technical French terms used in BAC exams
"""

import re
from pathlib import Path
from collections import Counter

# Technical vocabulary dictionary (French -> English)
TECHNICAL_VOCAB = {
    # Chemistry
    'cinétique': 'kinetics',
    'vitesse de réaction': 'reaction rate',
    'catalyseur': 'catalyst',
    'demi-vie': 'half-life',
    'oxydation': 'oxidation',
    'réduction': 'reduction',
    'redox': 'redox',
    'oxydant': 'oxidizing agent',
    'réducteur': 'reducing agent',
    'acide': 'acid',
    'base': 'base',
    'pH': 'pH',
    'pKa': 'pKa',
    'tampon': 'buffer',
    'titration': 'titration',
    'indicateur': 'indicator',
    'alcool': 'alcohol',
    'aldéhyde': 'aldehyde',
    'cétone': 'ketone',
    'ester': 'ester',
    'amine': 'amine',
    'acide carboxylique': 'carboxylic acid',
    'équilibre chimique': 'chemical equilibrium',
    'constante d\'équilibre': 'equilibrium constant',
    'Le Chatelier': 'Le Chatelier\'s principle',
    'réactif': 'reactant',
    'produit': 'product',
    'concentration': 'concentration',
    'mole': 'mole',
    'molaire': 'molar',
    
    # Physics - Mechanics
    'énergie mécanique': 'mechanical energy',
    'énergie cinétique': 'kinetic energy',
    'énergie potentielle': 'potential energy',
    'ressort': 'spring',
    'oscillation': 'oscillation',
    'projectile': 'projectile',
    'friction': 'friction',
    'force': 'force',
    'accélération': 'acceleration',
    'vitesse': 'velocity/speed',
    'déplacement': 'displacement',
    'travail': 'work',
    'puissance': 'power',
    
    # Physics - Electromagnetism
    'champ électrique': 'electric field',
    'champ magnétique': 'magnetic field',
    'induction': 'induction',
    'Lorentz': 'Lorentz force',
    'flux': 'flux',
    'bobine': 'coil',
    'auto-induction': 'self-induction',
    'f.e.m': 'EMF (electromotive force)',
    'inductance': 'inductance',
    'capacité': 'capacitance',
    'condensateur': 'capacitor',
    'résistance': 'resistance',
    'résistor': 'resistor',
    'courant': 'current',
    'tension': 'voltage',
    'charge': 'charge',
    'électron': 'electron',
    'proton': 'proton',
    
    # Physics - AC Circuits
    'RLC': 'RLC circuit',
    'résonance': 'resonance',
    'impédance': 'impedance',
    'déphasage': 'phase shift',
    'Fresnel': 'Fresnel diagram',
    'fréquence': 'frequency',
    'période': 'period',
    'amplitude': 'amplitude',
    'sinusoïdale': 'sinusoidal',
    
    # Physics - Waves
    'onde': 'wave',
    'longueur d\'onde': 'wavelength',
    'interférence': 'interference',
    'diffraction': 'diffraction',
    'réflexion': 'reflection',
    'réfraction': 'refraction',
    
    # General Lab Terms
    'expérience': 'experiment',
    'mesure': 'measurement',
    'courbe': 'curve/graph',
    'graphique': 'graph',
    'tableau': 'table',
    'résultat': 'result',
    'conclusion': 'conclusion',
    'hypothèse': 'hypothesis',
    'théorie': 'theory',
    'formule': 'formula',
    'équation': 'equation',
    'solution': 'solution',
    'réaction': 'reaction',
    'composé': 'compound',
    'élément': 'element',
    'atome': 'atom',
    'molécule': 'molecule',
    'ion': 'ion',
    'cristal': 'crystal',
    'solide': 'solid',
    'liquide': 'liquid',
    'gaz': 'gas',
    'température': 'temperature',
    'pression': 'pressure',
    'volume': 'volume',
    'masse': 'mass',
    'densité': 'density',
}

# Common French expressions in exams
COMMON_EXPRESSIONS = {
    'Définir': 'Define',
    'Calculer': 'Calculate',
    'Déterminer': 'Determine',
    'Montrer que': 'Show that',
    'En déduire': 'Deduce/Infer',
    'Justifier': 'Justify',
    'Expliquer': 'Explain',
    'Comparer': 'Compare',
    'Analyser': 'Analyze',
    'Interpréter': 'Interpret',
    'Écrire l\'équation': 'Write the equation',
    'Nommer': 'Name',
    'Identifier': 'Identify',
    'Classer': 'Classify',
    'Ordonner': 'Order',
    'Vérifier': 'Verify',
    'Estimer': 'Estimate',
    'Évaluer': 'Evaluate',
    'Prédire': 'Predict',
    'Proposer': 'Propose',
}

class FrenchVocabularyExtractor:
    def __init__(self, exams_dir):
        self.exams_dir = Path(exams_dir)
        self.vocab_found = Counter()
        self.expressions_found = Counter()
        
    def extract_from_all_exams(self):
        """Extract vocabulary from all exam files"""
        exam_files = list(self.exams_dir.glob('BAC-*.md'))
        
        for exam_file in exam_files:
            self._extract_from_file(exam_file)
        
        return self._compile_results()
    
    def _extract_from_file(self, exam_path):
        """Extract vocabulary from a single exam"""
        with open(exam_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Find technical vocabulary
        for french, english in TECHNICAL_VOCAB.items():
            pattern = r'\b' + re.escape(french) + r'\b'
            matches = len(re.findall(pattern, content, re.IGNORECASE))
            if matches > 0:
                self.vocab_found[french] += matches
        
        # Find common expressions
        for expression, translation in COMMON_EXPRESSIONS.items():
            pattern = r'\b' + re.escape(expression) + r'\b'
            matches = len(re.findall(pattern, content, re.IGNORECASE))
            if matches > 0:
                self.expressions_found[expression] += matches
    
    def _compile_results(self):
        """Compile extraction results"""
        return {
            'vocabulary': [
                {
                    'french': term,
                    'english': TECHNICAL_VOCAB[term],
                    'frequency': count,
                    'importance': 'HIGH' if count >= 10 else 'MEDIUM' if count >= 5 else 'LOW'
                }
                for term, count in self.vocab_found.most_common()
            ],
            'expressions': [
                {
                    'french': expr,
                    'english': COMMON_EXPRESSIONS[expr],
                    'frequency': count
                }
                for expr, count in self.expressions_found.most_common()
            ]
        }


def main():
    exams_dir = Path(__file__).parent / 'Study-Vault/04-Exams/BAC-2002-2012'
    
    if not exams_dir.exists():
        print(f"Error: Directory not found: {exams_dir}")
        return
    
    print("📚 Extracting French Vocabulary from BAC Exams...")
    print("=" * 70)
    
    extractor = FrenchVocabularyExtractor(exams_dir)
    results = extractor.extract_from_all_exams()
    
    # Print vocabulary
    print("\n🔤 TECHNICAL VOCABULARY (Most Frequent)")
    print("-" * 70)
    print(f"{'French':<30} {'English':<30} {'Freq':<8} {'Importance'}")
    print("-" * 70)
    
    for item in results['vocabulary'][:30]:
        print(f"{item['french']:<30} {item['english']:<30} {item['frequency']:<8} {item['importance']}")
    
    # Print expressions
    print("\n\n📝 COMMON EXAM EXPRESSIONS")
    print("-" * 70)
    print(f"{'French':<30} {'English':<30} {'Frequency'}")
    print("-" * 70)
    
    for item in results['expressions']:
        print(f"{item['french']:<30} {item['english']:<30} {item['frequency']}")
    
    # Save to file
    output_file = Path(__file__).parent / 'FRENCH-VOCABULARY-GUIDE.md'
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write("# 📚 French Vocabulary Guide for BAC Exams\n\n")
        f.write("## Technical Vocabulary\n\n")
        f.write("| French | English | Frequency | Importance |\n")
        f.write("|--------|---------|-----------|------------|\n")
        
        for item in results['vocabulary']:
            f.write(f"| {item['french']} | {item['english']} | {item['frequency']} | {item['importance']} |\n")
        
        f.write("\n## Common Exam Expressions\n\n")
        f.write("| French | English | Frequency |\n")
        f.write("|--------|---------|----------|\n")
        
        for item in results['expressions']:
            f.write(f"| {item['french']} | {item['english']} | {item['frequency']} |\n")
    
    print(f"\n💾 Vocabulary guide saved to: {output_file}")

if __name__ == '__main__':
    main()
