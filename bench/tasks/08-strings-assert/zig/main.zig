const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    const word = "hello";
    const length = word.len;
    const doubled = word ++ word;
    try stdout.print("{}\n", .{length});
    try stdout.print("{s}\n", .{doubled});
    std.debug.assert(length == 5);
}
