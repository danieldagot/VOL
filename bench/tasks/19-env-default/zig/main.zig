const std = @import("std");

pub fn main() !void {
    const stdout = std.io.getStdOut().writer();
    try std.posix.setenv("VOL_DENSITY_PORT", "9090", true);

    if (std.posix.getenv("VOL_DENSITY_PORT")) |v| {
        try stdout.print("{s}\n", .{v});
    } else {
        try stdout.print("missing\n", .{});
    }
    if (std.posix.getenv("VOL_DENSITY_NO_SUCH_KEY")) |v| {
        try stdout.print("{s}\n", .{v});
    } else {
        try stdout.print("fallback\n", .{});
    }
}
