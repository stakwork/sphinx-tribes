import { getBountyStatus } from '../filterTable';

describe('getBountyStatus', () => {
  test('should return Open status', () => {
    const result = getBountyStatus('open');
    expect(result.Open).toBe(true);
    expect(result.Assigned).toBe(false);
    expect(result.Paid).toBe(false);
  });

  test('should return In Progress status', () => {
    const result = getBountyStatus('in-progress');
    expect(result.Open).toBe(false);
    expect(result.Assigned).toBe(true);
    expect(result.Paid).toBe(false);
  });

  test('should return Completed status', () => {
    const result = getBountyStatus('completed');
    expect(result.Open).toBe(false);
    expect(result.Assigned).toBe(false);
    expect(result.Paid).toBe(true);
  });

  test('should return default status for unknown criterion', () => {
    const result = getBountyStatus('unknown' as any);
    expect(result.Open).toBe(false);
    expect(result.Assigned).toBe(false);
    expect(result.Paid).toBe(false);
  });
});
