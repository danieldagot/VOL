const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    const scores = [_]i64{ 85, 72, 91, 60, 78, 95, 55, 68, 88, 74 };
    var total: i64 = 0;
    var a_grades: i64 = 0;
    var b_grades: i64 = 0;
    var passing: i64 = 0;
    var failing: i64 = 0;
    for (scores) |s| {
        total += s;
        if (s >= 90) a_grades += 1;
        if (s >= 80 and s < 90) b_grades += 1;
        if (s >= 60) passing += 1;
        if (s < 60) failing += 1;
    }
    const avg = total / @as(i64, @intCast(scores.len));
    try stdout.print("Class average: {}\n", .{avg});
    try stdout.print("A grades: {}\n", .{a_grades});
    try stdout.print("B grades: {}\n", .{b_grades});
    try stdout.print("Passing: {}\n", .{passing});
    try stdout.print("Failing: {}\n", .{failing});
}
