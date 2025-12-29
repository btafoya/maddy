### [\#757] limits support for target.smtp
**URL**: https://github.com/foxcpp/maddy/issues/757
**State**: OPEN
**Labels**: new feature, good first issue, 

**Body:**
# Use case

Forwarding messages to external systems with ratelimits.

# Your idea for a solution

Add limits to target.smtp - the code will be very similar to target.remote.

As an alternative consider implementing target-agnostic "target.limited" wrapper. 

Consider interaction with target.queue - might be a good idea to skip increasing retries if delivery fails due to timeout waiting on limits.