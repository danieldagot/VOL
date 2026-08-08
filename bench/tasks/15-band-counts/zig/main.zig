const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    const vals = [_]i64{ 22, 18, 25, 31, 29, 17, 24, 28, 20, 26 };
    var total: i64 = 0;
    var hot: i64 = 0;
    var mild: i64 = 0;
    var cold: i64 = 0;
    for (vals) |v| {
        total += v;
        if (v >= 28) {
            hot += 1;
        } else if (v >= 20) {
            mild += 1;
        } else {
            cold += 1;
        }
    }
    try stdout.print("{}\n", .{vals.len});
    try stdout.print("{}\n", .{@divTrunc(total, @as(i64, @intCast(vals.len)))});
    try stdout.print("{}\n", .{hot});
    try stdout.print("{}\n", .{mild});
    try stdout.print("{}\n", .{cold});
}
