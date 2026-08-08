const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    const nums = [_]i64{ 3, 8, 1, 12, 5, 9, 4, 15, 2, 11 };
    var large_count: i64 = 0;
    var large_sum: i64 = 0;
    var large_double: i64 = 0;
    var small_count: i64 = 0;
    for (nums) |n| {
        if (n > 5) {
            large_count += 1;
            large_sum += n;
            large_double += n * 2;
        }
        if (n < 5) {
            small_count += 1;
        }
    }
    try stdout.print("{}\n", .{large_count});
    try stdout.print("{}\n", .{large_sum});
    try stdout.print("{}\n", .{large_double});
    try stdout.print("{}\n", .{small_count});
}
