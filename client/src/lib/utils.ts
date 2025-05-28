import { clsx, type ClassValue } from "clsx";
import { toast } from "sonner";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs));
}

export function isValidYouTubeWatchUrl(url: string) {
    const regex = /^https?:\/\/(www\.)?youtube\.com\/watch\?v=[\w-]{11}(&.*)?$/;
    return regex.test(url);
}

export function isValidTimestamp(time: string): boolean {
    const parts = time.trim().split(":");

    if (parts.length < 1 || parts.length > 3) return false;
    if (parts.some((p) => p === "" || isNaN(Number(p)))) return false;

    const nums = parts.map(Number);
    const [h, m, s] =
        parts.length === 3
            ? nums
            : parts.length === 2
                ? [0, ...nums]
                : [0, 0, ...nums];

    if (m >= 60 || s >= 60 || h < 0 || m < 0 || s < 0) return false;

    return true;
}

export function toTotalSeconds(time: string): number {
    const parts = time.trim().split(":").map(Number);

    if (parts.length === 3) return parts[0] * 3600 + parts[1] * 60 + parts[2];
    if (parts.length === 2) return parts[0] * 60 + parts[1];
    if (parts.length === 1) return parts[0];
    return 0;
}

export function validateInputs({
    url,
    startTime,
    endTime,
}: {
    url: string;
    startTime: string;
    endTime: string;
}) {
    if (!isValidYouTubeWatchUrl(url)) {
        toast.error("Please enter a valid YouTube URL.");
        return false;
    }

    if (!isValidTimestamp(startTime) || !isValidTimestamp(endTime)) {
        toast.error("Please enter a valid time in HH:MM:SS, MM:SS, or SS format.");
        return false;
    }

    const start = toTotalSeconds(startTime);
    const end = toTotalSeconds(endTime);

    if (start >= end) {
        toast.error("Start time must be before end time.");
        return false;
    }

    return true;
}
