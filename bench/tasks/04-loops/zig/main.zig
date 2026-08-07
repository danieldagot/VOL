const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    try stdout.print("Countdown\n", .{});
    var remaining: i64 = 3;
    while (remaining > 0) {
        try stdout.print("{}\n", .{remaining});
        remaining -= 1;
    }
    for (0..2) |_| {
        try stdout.print("Go!\n", .{});
    }
}
