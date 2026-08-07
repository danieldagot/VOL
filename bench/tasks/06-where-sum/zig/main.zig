const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    const numbers = [_]i64{ 4, 7, 2, 9, 12 };
    var total: i64 = 0;
    for (numbers) |n| {
        if (n > 5) total += n;
    }
    try stdout.print("Sum: {}\n", .{total});
    std.debug.assert(total == 28);
}
