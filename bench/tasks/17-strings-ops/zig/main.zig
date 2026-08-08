const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    var gpa = std.heap.page_allocator;

    const trimmed = std.mem.trim(u8, "  vol  ", " \t\n\r");
    try stdout.print("{s}\n", .{trimmed});

    var it = std.mem.splitScalar(u8, "a,b,c", ',');
    var parts = std.ArrayList([]const u8).init(gpa);
    defer parts.deinit();
    while (it.next()) |part| {
        try parts.append(part);
    }
    const joined = try std.mem.join(gpa, "-", parts.items);
    defer gpa.free(joined);
    try stdout.print("{s}\n", .{joined});

    try stdout.print("{}\n", .{std.mem.indexOf(u8, "vocabulary", "cab") != null});
    const replaced = try std.mem.replaceOwned(u8, gpa, "a-a-a", "-", "+");
    defer gpa.free(replaced);
    try stdout.print("{s}\n", .{replaced});
}
