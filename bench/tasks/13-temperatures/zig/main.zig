const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    const temps = [_]i64{ 22, 18, 25, 31, 29, 17, 24, 28, 20, 26 };
    var total: i64 = 0;
    var hot: i64 = 0;
    var mild: i64 = 0;
    var cold: i64 = 0;
    for (temps) |t| {
        total += t;
        if (t >= 28) {
            hot += 1;
        } else if (t >= 20) {
            mild += 1;
        } else {
            cold += 1;
        }
    }
    const avg = total / @as(i64, @intCast(temps.len));
    try stdout.print("Days measured: {}\n", .{temps.len});
    try stdout.print("Average: {}\n", .{avg});
    try stdout.print("Hot days (28+): {}\n", .{hot});
    try stdout.print("Mild days: {}\n", .{mild});
    try stdout.print("Cold days (<20): {}\n", .{cold});
}
