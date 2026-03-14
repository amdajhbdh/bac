#!/usr/bin/env nu

mkdir /home/med/Downloads/bac-ocr/mauritania
mkdir /home/med/Downloads/bac-ocr/mauritania-extra
mkdir /home/med/Downloads/bac-ocr/pc-sm

let maurPdfs = glob /home/med/Downloads/mauritania/*.pdf
let extraPdfs = glob /home/med/Downloads/mauritania-extra/*.pdf
let pcPdfs = glob /home/med/Downloads/pc-sm/*.pdf

print $"Mauritania: ( $maurPdfs | length ) files"
print $"Extra: ( $extraPdfs | length ) files"
print $"PC: ( $pcPdfs | length ) files"

for pdf in $maurPdfs {
    let fname = ( $pdf | path parse | get stem )
    let outFile = $"/home/med/Downloads/bac-ocr/mauritania/($fname).txt"
    if ( not ( $outFile | path exists ) ) {
        print $fname
        let tmp = $"/tmp/ocr_($fname)"
        ^pdftoppm -png -r 150 $pdf $tmp
        for img in (glob $"($tmp)-*.png") {
            ^tesseract $img stdout -l fra+eng | save -a $outFile
            rm $img
        }
    }
}

for pdf in $extraPdfs {
    let fname = ( $pdf | path parse | get stem )
    let outFile = $"/home/med/Downloads/bac-ocr/mauritania-extra/($fname).txt"
    if ( not ( $outFile | path exists ) ) {
        print $fname
        let tmp = $"/tmp/ocr_($fname)"
        ^pdftoppm -png -r 150 $pdf $tmp
        for img in (glob $"($tmp)-*.png") {
            ^tesseract $img stdout -l fra+eng | save -a $outFile
            rm $img
        }
    }
}

for pdf in $pcPdfs {
    let fname = ( $pdf | path parse | get stem )
    let outFile = $"/home/med/Downloads/bac-ocr/pc-sm/($fname).txt"
    if ( not ( $outFile | path exists ) ) {
        print $fname
        let tmp = $"/tmp/ocr_($fname)"
        ^pdftoppm -png -r 150 $pdf $tmp
        for img in (glob $"($tmp)-*.png") {
            ^tesseract $img stdout -l fra+eng | save -a $outFile
            rm $img
        }
    }
}

print "=== Done ==="
