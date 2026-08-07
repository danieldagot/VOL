const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    const player1 = [_]i64{ 8, 6, 9, 5, 10, 7 };
    const player2 = [_]i64{ 7, 8, 6, 9, 8, 8 };
    var p1_total: i64 = 0;
    var p2_total: i64 = 0;
    for (player1) |s| p1_total += s;
    for (player2) |s| p2_total += s;
    const winner = if (p1_total > p2_total) "Player 1" else "Player 2";
    var p1_strong: i64 = 0;
    var p2_strong: i64 = 0;
    for (player1) |s| { if (s >= 8) p1_strong += 1; }
    for (player2) |s| { if (s >= 8) p2_strong += 1; }
    try stdout.print("Player 1 total: {}\n", .{p1_total});
    try stdout.print("Player 2 total: {}\n", .{p2_total});
    try stdout.print("Winner: {s}\n", .{winner});
    try stdout.print("P1 rounds 8+: {}\n", .{p1_strong});
    try stdout.print("P2 rounds 8+: {}\n", .{p2_strong});
}
