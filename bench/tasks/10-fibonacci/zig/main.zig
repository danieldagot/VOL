const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    var a: i64 = 0;
    var b: i64 = 1;
    for (0..8) |_| {
        try stdout.print("{}\n", .{a});
        const c = a + b;
        a = b;
        b = c;
    }
}
