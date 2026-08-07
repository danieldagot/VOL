const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    const a = 10;
    const b = 3;
    const sum = a + b;
    const product = a * b;
    const diff = a - b;
    try stdout.print("{}\n", .{sum});
    try stdout.print("{}\n", .{product});
    try stdout.print("{}\n", .{diff});
}
