import { expect, test } from 'vitest'
import { formatCentimeters, calculateCentimeters, calculateMilliseconds } from ".";

test('formatCentimeters ftin', () => {
    expect(formatCentimeters(2.54, "ftin")).toBe("1\"")
    expect(formatCentimeters(30.5, "ftin")).toBe("1\' 0.01\"")
    expect(formatCentimeters(30.48, "ftin")).toBe("1\'")
    expect(formatCentimeters(162.56, "ftin")).toBe("5\' 4\"")
    expect(formatCentimeters(152.4, "ftin")).toBe("5\'")
    expect(formatCentimeters(91.44, "ftin")).toBe("3\'")
})

test('formatCentimeters cm', () => {
    expect(formatCentimeters(2.55, "cm")).toBe("2.55 cm")
})

test('calculateCentimeters ftin', () => {
    expect(calculateCentimeters(5, 4, "ftin")).toBe(162.56)
    expect(calculateCentimeters(5, 0, "ftin")).toBe(152.4)
})

test('calculateMilliseconds', () => {
    expect(calculateMilliseconds("1:00")).toBe(60000)
    expect(calculateMilliseconds("1:00.00")).toBe(60000)
    expect(calculateMilliseconds("1:00.001")).toBe(60001)
    expect(calculateMilliseconds("1:00.01")).toBe(60010)
    expect(calculateMilliseconds("1:00.10")).toBe(60100)
    expect(calculateMilliseconds("1:10.00")).toBe(70000)
    expect(calculateMilliseconds("1:11.00")).toBe(71000)
})