import { createPinia, setActivePinia } from "pinia";
import { useResultsStore } from "./results";

describe("Results Store", () => {
	beforeEach(() => {
		setActivePinia(createPinia());
	});

	it("recordFinishTime", () => {
		const sut = useResultsStore();

		sut.recordFinishTime(finishTime1);
		expect(sut.finishers).toBe(1);
		sut.recordFinishTime(finishTime2);
		expect(sut.finishers).toBe(2);
	});

	it("recordRunnerResult", () => {
		const sut = useResultsStore();

		sut.recordRunnerResult(finisher1);
		expect(sut.finishers).toBe(1);
		sut.recordRunnerResult(finisher2);
		expect(sut.finishers).toBe(2);
	});

	it("recordFinishTime -> recordRunnerResult", () => {
		const sut = useResultsStore();

		sut.recordFinishTime(finishTime1);
		expect(sut.finishers).toBe(1);
		sut.recordRunnerResult(finisher1);
		expect(sut.finishers).toBe(1);
	});

	it("recordRunnerResult -> recordFinishTime", () => {
		const sut = useResultsStore();
		sut.recordRunnerResult(finisher1);
		expect(sut.finishers).toBe(1);
		sut.recordFinishTime(finishTime1);
		expect(sut.finishers).toBe(2);
	});
});

const finishTime1 = {
	timestamp: Date.UTC(1983, 6, 13, 23, 43, 34, 1),
	minutes: 14,
	seconds: 23,
	milliseconds: 23,
};

const finishTime2 = {
	timestamp: Date.UTC(1983, 6, 14, 23, 43, 34, 1),
	minutes: 15,
	seconds: 23,
	milliseconds: 23,
};

const finisher1 = {
	timestamp: Date.UTC(1983, 6, 13, 23, 43, 34, 1),
	bib: 121,
};

const finisher2 = {
	timestamp: Date.UTC(1983, 6, 13, 23, 43, 35, 1),
	bib: 231,
};
