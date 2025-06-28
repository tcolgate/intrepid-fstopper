%.pdf: %.dat
	pdflatex --jobname "$(shell echo $(<) | cut -d. -f1)" curve.tex
