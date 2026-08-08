const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    const xs = [_]i64{ 1, 2, 3, 4, 5, 6, 7, 8 };
    var mid_count: i64 = 0;
    var high_sum: i64 = 0;
    for (xs) |x| {
        const n = x * 3;
        if (n > 10 and n < 20) {
            mid_count += 1;
        }
        if (n > 10) {
            high_sum += n;
        }
    }
    try stdout.print("{}\n", .{mid_count});
    try stdout.print("{}\n", .{high_sum});
}
