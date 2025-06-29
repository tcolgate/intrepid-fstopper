FILES=\
  IlfordFP4+-FX55\ 50+10+440-20c-8m30s-Ilford\ Agit.pdf \
  IlfordFP4+-Pyrocat\ HDC\ 1+1+100-20c-13m-Ilford\ Agit.pdf \
  IlfordHP5+-FX55\ 50+10+440-20c-13m30s-Ilford\ Agit.pdf \
  IlfordHP5+-Pyrocat\ HDC\ 1+1+100-20c-15m-Ilford\ Agit.pdf \
  Kentmere200-Pyrocat\ HDC\ 1+1+100-20c-9m30s-Ilford\ Agit.pdf


all: $(FILES)

clean:
	rm -f $(FILES)

%.pdf: %.dat curve.tex
	pdflatex --jobname "$(shell echo $(<) | cut -d. -f1)" curve.tex
