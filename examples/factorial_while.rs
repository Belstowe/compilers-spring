fn factorial(n: i64) -> i64 {
    let result: i64 = 1;
    let tmp: i64 = n;
    while tmp > 1 {
        result = result * tmp;
        tmp = tmp - 1;
    }
    result
}

fn main() {
    ruster::writeln_i64(factorial(5));
}