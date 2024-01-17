import { BountyMetrics } from 'store/main';

export const normalizeMetrics = (data: any): BountyMetrics => ({
  bounties_posted: data.BountiesPosted || data.bounties_posted,
  bounties_paid: data.BountiesPaid || data.bounties_paid,
  bounties_paid_average: data.bounties_paid_average || data.BountiesPaidPercentage,
  sats_posted: data.sats_posted || data.SatsPosted,
  sats_paid: data.sats_paid || data.SatsPaid,
  sats_paid_percentage: data.sats_paid_percentage || data.SatsPaidPercentage,
  average_paid: data.average_paid || data.AveragePaid,
  average_completed: data.average_completed || data.AverageCompleted,
  unique_hunters_paid: data.unique_hunters_paid || data.uniqueHuntersPaid,
  new_hunters_paid: data.new_hunters_paid || data.newHuntersPaid
});
