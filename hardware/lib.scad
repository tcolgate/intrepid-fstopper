use <lib.scad>

module trayOuter(
  w,
  d,
  h,
  cornerRadius=5,
) {
    translate([(w * -0.5),(d * -0.5),0]){
      minkowski(){
        cube([w,d,h]);
        cylinder(r=cornerRadius,h=1);
    }
  }
}

module trayInner(
  w,
  h,
  cornerRadius=2,
) {
    translate([(w * -0.5),(w * -0.5),0]){
      cube([w,w,h]);
       
      minkowski(){
        cube([w,w,h]);
        cylinder(r=cornerRadius,h=1);
    }
  }
}

module trayInsert(
  w,
  h,
  filterSize,
  cornerRadius=2,
) {
    difference(){
      trayInner(w,h,cornerRadius);
      cylinder(h=10,r=(filterSize/2 - 1),center=true);
    };
};


module filterCarrier(
  width, 
  depth, 
  height, 
  filterSize
) {
  union(){
      // tray
      difference(){
        difference(){
          // tray outer
          trayOuter(w=width,d=depth,h=height);
          // tray inner void (1.5 is the thickness of the bottom plate)
          translate([0,0,1.5]){
            trayInner(w=filterSize,h=h);    
          }
        };
        // central hole
        cylinder(h=10,r=(filterSize/2 - 1),center=true);
      }

      // handle
      translate([(-width/2 -10),0,1]){
          union(){
              //stem
              cube([11,10,2],center=true);
              //handle
              translate([-5,0,1.5]){
                  cube([2.5,13,5],center=true);
              }
          }
      }
  }
};


filterCarrier(
  width=52.5,
  depth=52.5,
  height=2.5,
  filterSize=50
);

translate([60,0,0]){
  trayInsert(
    w=50,
    h=1.5,
    filterSize=50
  );
};
