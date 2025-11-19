# Useful modelss for the Intrepid Compact Enlarger

There are two models here:
* an Under the Lens filter tray for Ilford Multigrade filters
* a clip-together gel filter holder similar to the Ilford holder. Intended
  for holding Kodak wratten gel filters for process like Tri-Colour printing.

I also highly recommend [Durst Nepla lens board upgrade](https://bluegrassphotographics.com/Lensboard-Upgrade-for-Intrepid-MF-Enlarger-p719001905) from Gregory Davis.
You'll need at least one of his Nepla lens holders too. I find the combination of his board and my filter tray makes for a really nice setup. I did add an ring of 1mm light seal foam around the underside of his lens holder (the side that faces the enlarger), it's not crucial, but it stopped my lens slipping when I change aperture or use the aperture lever.

## Under the lens filter tray (`carrier.scad`)

The model assume [this exact clamp](https://www.amazon.co.uk/dp/B09NFDXC25?ref_=ppx_hzsearch_conn_dt_b_fed_asin_title_2&th=1), a "Suxing Sleeve Swivel Clamp Chuck for Magnetic Stands Holder Bar Dial Indicator Gauge (D6/D8-D8)"
. I ordered these off amazon. It affixes to the 8mm rod that carries the
enlarger lens board, and the model has an 8mm peg to fit it to the clamp. If
you use an alternative you will need to adjust the parameters of the model so
that the peg at the read fits the hole on your clamp, and that the peg is
correctly position relative to the centre hole so that the filter hole lines up
precisely under the lens.

This holder fits the under-the-lens filters from the Ilford kit.

I am planning to create a second version of this that will also hold Agfa
colour filters for RA4 printing.

### Making the original Ilford holder fit

I had a hard time fitting the Ilford Multigrade Under the Lens multigrade
filter kit with the Compact Enlarger.

I use Schneider enlarger lenses which have a short thread and an extra arm to
open the aperture. The short thread means that the lens will not screw in properly
with the 3-legged holder between it and the lens board. Also, it becomes very
awkward to change lenses and keep the holder in place. In addition, once fitted
the three legged holder gets in the way of the preview arm. The approach below seemed
to alleviate a lot of these problems, but was still annoying to use.

I did manage to make it fit by:

* File out the inner circle of the 3-legged holder so that the lens will fit
  entirely inside it without "clamping" it to the lens board.
* The the legs on the three-legged holder are arrange with two of the
  legscloser together. The section with the two closer together faces "forward"
  to the entrance of the filter holder.
* Assume the filter holder will fit "side ways", that is, rather than the filters going
  in the front facing you, they go in from the left hand side, so we'll arrange for the
  two closer legs to be to the left of the lens hole. This makes it easier to use the
  aperturen opening preview arm.
* Cut out the long section that will be  toward the back of the lens board,
  this makes it easier to fit the 3-legged holder when the lens is already
  mounted to the lens board. The gap makes room for 3 legged arm to slide on without
  hitting the preview arm.
* Use some small bits of sticky velcro to let me affix the 3-legged holder around
  the lens hole on the board. You don't need loads, the holder is not heavy. Make sure
  the velcro does not get between the lens board and the lens.

## 50mm Gel Filter holder

This is split between two files `filter-holder-50mm.scad`,
`filter-holder-50mm-insert.scad`. The first is the main body of the holder. The
second is a push in insert. You will need to sand or file down the small detents
on the edge until you get a comfortable insert into the holder body.
