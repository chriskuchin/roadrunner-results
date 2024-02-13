import { expect, test } from 'vitest'
import { formatCentimeters } from ".";

test('formatCentimeters', () => {
    expect(formatCentimeters(2.54, "ftin")).toBe("1\"")
    expect(formatCentimeters(30.5, "ftin")).toBe("1\' 0.01\"")
    expect(formatCentimeters(30.48, "ftin")).toBe("1\'")
    expect(formatCentimeters(162.56, "ftin")).toBe("5\' 4\"")
})
