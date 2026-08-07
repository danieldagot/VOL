const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    const language = "VOL";
    const version = 1;
    try stdout.print("Hello from {s}\n", .{language});
    try stdout.print("Prototype version {}\n", .{version});
}
