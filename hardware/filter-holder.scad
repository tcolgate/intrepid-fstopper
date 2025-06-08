module filterCarrier(
) {
  union(){
      // tray
      difference(){
        difference(){
          // tray outer
          minkowski()
          {
            cube([52.5,52.5,5.5]);
            cylinder(r=5,h=1);
          }
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
        translate([26.25,26.25,-5]){
          cylinder(h=10,r=(51.8/2));
        };
      }

      // handle
      translate([-12,((52.5 - 13)/2),0]){
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
}

filterCarrier(){
}
