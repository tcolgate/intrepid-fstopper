// This is a test for a push fit onto the
// enlarger guide post
difference(){
    cylinder(5,16,16, centre=true);
    translate([0,0,-4]){
            cylinder(10,8,8, centre=true);
    }
}
