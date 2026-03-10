#!/usr/bin/env python3
"""
BAC PDF Downloader & Organizer
Downloads publicly available BAC exam PDFs and organizes them
"""

import os
import requests
from pathlib import Path
from urllib.parse import urljoin
import json

# Direct URLs that are publicly available (no authentication needed)
BAC_URLS = {
    "Mathématiques": {
        "2025": {
            "C": "https://maurimath.net/documents/sbac/BacC2025SN.pdf",
            "D": "https://maurimath.net/documents/sbac/BacD2025SN.pdf",
        },
        "2024": {
            "C": "https://maurimath.net/documents/sbac/SujetBacC2024SN.pdf",
            "D": "https://maurimath.net/documents/sbac/SujetBacD2024SN.pdf",
        },
        "2022": {
            "C_Exp_SN": "https://maurimath.net/documents/sbac/BacCExpSN2022.pdf",
            "C_Exp_SC": "https://maurimath.net/documents/sbac/BacCExpSC2022.pdf",
            "D_Exp_SC": "https://maurimath.net/documents/sbac/BacDExpSC2022.pdf",
        },
        "2021": {
            "C": "https://maurimath.net/documents/sbac/BacC2021SN.pdf",
        },
        "2020": {
            "C": "https://maurimath.net/documents/sbac/BacC2020sn.pdf",
        },
        "Collections": {
            "2000-2012_C_TMGM": "https://www.sigmaths.net/bac2/mauritanie/baccmaths.pdf",
            "2022_M_TMGM": "https://www.sigmaths.net/bac2/mauritanie/M_TMGM2022_sujet_correction.pdf",
        }
    },
    "Physique-Chimie": {
        "Collections": {
            "BAC_D_2005-2020": "https://prepabac.app/sc/dwnld/bac_d_pc.pdf",
        }
    },
    "SVT": {
        "Collections": {
            "BAC_C_2010-2018": "https://prepabac.app/sc/dwnld/bac_c_sn.pdf",
        }
    }
}

class BACDownloader:
    def __init__(self, output_dir):
        self.output_dir = Path(output_dir)
        self.output_dir.mkdir(parents=True, exist_ok=True)
        self.session = requests.Session()
        self.session.headers.update({
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
        })
        self.results = {
            'downloaded': [],
            'failed': [],
            'skipped': []
        }
    
    def download_file(self, url, filepath):
        """Download a single file"""
        try:
            print(f"  Downloading: {url}")
            response = self.session.get(url, timeout=10, allow_redirects=True)
            response.raise_for_status()
            
            with open(filepath, 'wb') as f:
                f.write(response.content)
            
            size_mb = filepath.stat().st_size / (1024*1024)
            print(f"  ✅ Saved: {filepath.name} ({size_mb:.1f} MB)")
            return True
        except Exception as e:
            print(f"  ❌ Failed: {str(e)}")
            return False
    
    def download_all(self):
        """Download all available PDFs"""
        print("🔍 BAC PDF Downloader")
        print("=" * 60)
        
        for subject, years in BAC_URLS.items():
            subject_dir = self.output_dir / subject
            subject_dir.mkdir(exist_ok=True)
            
            print(f"\n📚 {subject}")
            print("-" * 60)
            
            for year, series in years.items():
                year_dir = subject_dir / year
                year_dir.mkdir(exist_ok=True)
                
                for series_name, url in series.items():
                    filename = f"{series_name}.pdf"
                    filepath = year_dir / filename
                    
                    if filepath.exists():
                        print(f"  ⏭️  Skipped (exists): {filename}")
                        self.results['skipped'].append(str(filepath))
                    else:
                        if self.download_file(url, filepath):
                            self.results['downloaded'].append(str(filepath))
                        else:
                            self.results['failed'].append(f"{subject}/{year}/{series_name}")
        
        self._print_summary()
    
    def _print_summary(self):
        """Print download summary"""
        print("\n" + "=" * 60)
        print("📊 DOWNLOAD SUMMARY")
        print("=" * 60)
        print(f"✅ Downloaded: {len(self.results['downloaded'])}")
        print(f"❌ Failed: {len(self.results['failed'])}")
        print(f"⏭️  Skipped: {len(self.results['skipped'])}")
        
        if self.results['failed']:
            print("\n⚠️  Failed downloads:")
            for item in self.results['failed']:
                print(f"  - {item}")
        
        # Save manifest
        manifest = {
            'downloaded': self.results['downloaded'],
            'failed': self.results['failed'],
            'skipped': self.results['skipped'],
            'total': len(self.results['downloaded']) + len(self.results['failed']) + len(self.results['skipped'])
        }
        
        manifest_file = self.output_dir / 'download-manifest.json'
        with open(manifest_file, 'w') as f:
            json.dump(manifest, f, indent=2)
        
        print(f"\n💾 Manifest saved: {manifest_file}")

def main():
    output_dir = Path(__file__).parent / "BAC-Exams-PDFs"
    downloader = BACDownloader(output_dir)
    downloader.download_all()

if __name__ == '__main__':
    main()
