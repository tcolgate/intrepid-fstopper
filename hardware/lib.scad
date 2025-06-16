use <lib.scad>

module trayOuter(
  width,
  depth,
  height=3,
  cornerRadius=5,
) {
    translate([((width-(2*cornerRadius)) * -0.5),((depth-(2*cornerRadius)) * -0.5),0]){
      minkowski(){
        cube([(width-(2*cornerRadius)),(depth-(2*cornerRadius)),(height/2)]);
        cylinder(r=cornerRadius,h=(height/2));
    }
  }
}

module trayInner(
  width,
  height=1.5,
  cornerRadius=2,
) {
    translate([((width-(2*cornerRadius)) * -0.5),((width-(2*cornerRadius)) * -0.5),0]){
      minkowski(){
        cube([(width-(2*cornerRadius)),(width-(2*cornerRadius)),(height/2)]);
        cylinder(r=cornerRadius,h=(height/2));
      }
    }
}

module trayInsert(
  width,
  height=1.5,
  filterSize,
  cornerRadius=2,
  fudge=1,
) {
    difference(){
      union(){       
        trayInner(width,height,cornerRadius);
        for (i = [0, 90, 180, 270]){
          rotate([0, 0, i]){
            translate([(width/2 - 1),0,(height/2)]){
              cylinder(h=height,r=2,center=true);
            };
          };
        };            
      };
      cylinder(h=(height*2.2),r=(filterSize/2 - 1),center=true);
    };
};


module filterHolder(
  width,
  depth,
  height=3,
  base=1.5,
  filterSize
) {
  union(){
      // tray
      difference(){
        difference(){
          // tray outer
          trayOuter(width=width,depth=depth,height=height);
          // tray inner void
          #translate([0,0,(height - base)]){
            trayInner(width=filterSize,height=(height+1),cornerRadius=0);
          };
        };
        // central hole
        cylinder(h=10,r=(filterSize/2 - 1),center=true);
      }

      // handle
      translate([(-width/2 - 5),0,1]){
          union(){
              //stem
              cube([11,10,2],center=true);
              //handle
              translate([-5,0,base]){
                  cube([2.5,13,5],center=true);
              }
          }
      }
  }
};


filterHolder(
  width=63.8,
  depth=63.8,
  height=3,
  filterSize=50.3
);

translate([75,0,0]){
  trayInsert(
    width=50.3,
    height=1.5,
    filterSize=50.3
  );
};
