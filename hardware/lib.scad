use <lib.scad>

// Filter Holders with a retaining insert
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
              translate([-5,0,(2-(height/2))]){
                  cube([2.5,13,height],center=true);
              }
          }
      }
  }
};

// Ilford Filter under lens carrier

module attachmentPeg (
  pegWidth = 8,
  pegHeight = 16,
  pegBridgeLen, // This is the length from the back edge of carrier to the front flat of the peg
){
      union(){
         translate([-pegWidth/2,0,pegHeight/2]){
           rotate([0,90,0]){
             difference(){
               union(){
                   translate([pegHeight/2-pegHeight/4,0,0]){
                     cube([pegHeight/4,(pegWidth),(pegWidth)]);
                   }
                   translate([pegHeight/2-pegHeight/4,pegWidth/2,0]){
                     cube([(pegHeight/4),pegBridgeLen,(pegWidth)]);
                   };
               }
               translate([-pegHeight/8,0,-0.5]){
                 cylinder(r=((pegHeight+0.5)/2),h=(pegWidth+1));
               };
             };
           };
         };
         translate([0,0,pegHeight/2]){
         cylinder(r=(pegWidth/2),h=pegHeight,center=true);         }
      };
}

module mainBody (
  width = 66.3,
  depth = 69.3,
  wall = 2,
  rearGap = 4,
  aperture = 57.4,
  pegWidth = 7.5,
  pegHeight = 16.5,
  pegOffsetX = 18,
  pegOffsetY = 45,
){
   $forwardShift = (((aperture - depth)/2)+rearGap);
   $backEdgeLine = aperture/2 + rearGap;
   $pegBridgeLen = ((pegOffsetY - (pegWidth/2)) - $backEdgeLine);
   union(){
      translate([(pegOffsetX),(pegOffsetY),wall/2]){
        rotate([180,0,0]){
          attachmentPeg (
            pegWidth = pegWidth,
            pegHeight = pegHeight,
            pegBridgeLen = $pegBridgeLen
          );
        };
        #cylinder(r=1,h=10,center=true);
      };
      difference(){
          translate([0,$forwardShift,0]){
              cube([width,depth,wall],center=true);
          };
          cylinder(h=(wall+2),r=(aperture/2),center=true);
      }
  }
};



*filterHolder(
  width=63.8,
  depth=63.8,
  height=3,
  filterSize=50.3
);

*translate([75,0,0]){
  trayInsert(
    width=50.3,
    height=1.5,
    filterSize=50.3
  );
};

mainBody();

*attachmentPeg (
  pegWidth = 7.5,
  pegHeight = 16,
  pegBridgeLen= 10 // This is the length from the back edge of carrier to the front flat of the pe
);
