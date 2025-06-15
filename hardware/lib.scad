use <lib.scad>

module trayOuter(
  w,d,h
) {
  minkowski()
  {
    cube([w,d,h]);
    cylinder(r=5,h=1);
  }
}

module filterCarrier(
  width, depth, height, apertureRadius
) {
  union(){
      // tray
      difference(){
        difference(){
          // tray outer
          trayOuter(width,depth,height)
          // tray inner void
          translate([3.1,3.1,2.6]){
            minkowski()
            {
              cube([46.3,46.3,5.5]);
              cylinder(r=5,h=1);
            }
          }
        };
        // central hole
        translate([(width/2),(depth/2),-5]){
          cylinder(h=10,r=(51.8/2));
        };
      }

      // handle
      translate([-12,((width - 13)/2),0]){
          union(){
              //stem
              cube([11,10,2.5]);
              //handle
              translate([-2.5,-1.5,0]){
                  cube([2.5,13,6]);
              }
          }
      }
  }
};

