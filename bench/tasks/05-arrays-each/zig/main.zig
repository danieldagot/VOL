const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    var scores = [_]i64{ 72, 95, 81, 64 };
    try stdout.print("Students: {}\n", .{scores.len});
    scores[3] = 70;
    for (scores) |score| {
        if (score >= 80) {
            try stdout.print("High score: {}\n", .{score});
        }
    }
}
