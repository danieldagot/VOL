const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    const temperature = 24;
    const sunny = true;
    if (temperature >= 20 and sunny) {
        try stdout.print("Good weather for a walk.\n", .{});
    } else {
        try stdout.print("Stay inside and write VOL.\n", .{});
    }
}
