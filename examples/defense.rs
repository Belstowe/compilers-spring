fn a() {
    let b: i64 = 0;
    let c: i64 = 0;
    let d: i64;
    while b < 5 {
        c = c + 1;
        b = b + 1;
    }
    if c == 5 {
        d = 10;
    } else {
        d = 20;
    }
    let f: i64 = b + c + d;
    ruster::writeln_i64(b);
    ruster::writeln_i64(c);
    ruster::writeln_i64(d);
    ruster::writeln_i64(f);
}

fn main() {
    a();
}
