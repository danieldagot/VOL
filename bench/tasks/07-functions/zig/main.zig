const std = @import("std");

fn square(number: i64) i64 {
    return number * number;
}

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    try stdout.print("Hello, {s}\n", .{"friend"});
    try stdout.print("Six squared is {}\n", .{square(6)});
}
