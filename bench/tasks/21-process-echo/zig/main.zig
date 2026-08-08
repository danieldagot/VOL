const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    var gpa = std.heap.page_allocator;

    const result = try std.process.Child.run(.{
        .allocator = gpa,
        .argv = &[_][]const u8{ "echo", "vol" },
    });
    defer gpa.free(result.stdout);
    defer gpa.free(result.stderr);

    const status: i32 = switch (result.term) {
        .Exited => |code| @intCast(code),
        else => 1,
    };
    try stdout.print("{d}\n", .{status});
    try stdout.print("{s}\n", .{std.mem.trim(u8, result.stdout, " \t\n\r")});
}
