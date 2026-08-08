const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    var gpa = std.heap.page_allocator;

    const parsed = try std.json.parseFromSlice(std.json.Value, gpa, "{\"n\":3,\"name\":\"vol\"}", .{});
    defer parsed.deinit();
    const obj = parsed.value.object;
    try stdout.print("{d}\n", .{obj.get("n").?.integer});
    try stdout.print("{s}\n", .{obj.get("name").?.string});

    var body = std.json.ObjectMap.init(gpa);
    defer body.deinit();
    try body.put("n", .{ .integer = 3 });
    var out = std.ArrayList(u8).init(gpa);
    defer out.deinit();
    try std.json.stringify(std.json.Value{ .object = body }, .{}, out.writer());
    try stdout.print("{s}\n", .{out.items});
}
