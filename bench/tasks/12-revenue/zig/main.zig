const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    const revenue = [_]i64{ 240, 175, 289, 150, 225, 199, 180, 178 };
    var total: i64 = 0;
    var high_sum: i64 = 0;
    var high_count: i64 = 0;
    var budget_count: i64 = 0;
    for (revenue) |r| {
        total += r;
        if (r >= 200) {
            high_count += 1;
            high_sum += r;
        } else {
            budget_count += 1;
        }
    }
    try stdout.print("Total revenue: {}\n", .{total});
    try stdout.print("Premium orders (200+): {}\n", .{high_count});
    try stdout.print("Premium revenue: {}\n", .{high_sum});
    try stdout.print("Budget orders: {}\n", .{budget_count});
}
