const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    var gpa = std.heap.page_allocator;

    const joined = try std.fs.path.join(gpa, &[_][]const u8{ "reports", "2026", "q1.json" });
    defer gpa.free(joined);
    try stdout.print("{s}\n", .{joined});

    const path = "reports/2026/q1.json";
    try stdout.print("{s}\n", .{std.fs.path.basename(path)});
    try stdout.print("{s}\n", .{std.fs.path.dirname(path).?});
    try stdout.print("{s}\n", .{std.fs.path.extension(path)});
}
