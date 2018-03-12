# Differential
[![forthebadge](https://forthebadge.com/images/badges/you-didnt-ask-for-this.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/fo-shizzle.svg)](https://forthebadge.com)

![](https://cdn-images-1.medium.com/max/1600/1*ktMRyTnut5YK_0LoEJebUA.gif)

## Overview
This is the first pass at a [Differential Privacy](https://en.wikipedia.org/wiki/Differential_privacy) extension for PostgresSQL. 


### Installation
Follow the steps in the Dockerfile. You'll need to get `plgo` and use that to create the Makefile for the extension. Run `sudo make install` on the box running Postgres 9.6 and then run the `CREATE EXTENSION` command after rebooting the database.

### Example
Let's start with a contrived example, a dating website. The following is an example `users` table from our fictional dating app.The reason I say this is contrived, is because its unlikely it makes sense to denormalize our data like this. Its more likely that we have some kind of generic attribute table that contains relationships between users and attributes. But then again I've seen this table setup in more than one company I've worked for so ¯\_(ツ)_/¯.


| user_id | is_420_friendly | is_veteran | is_smoker |
|---------|-----------------|------------|-----------|
| 1       | T               | F          | F         |
| 2       | F               | T          | T         |
| 3       | F               | T          | F         |
| 4       | T               | F          | T         |



Anyway, each row corresponds to a user. For each user, we have a boolean field representing if they ingest Marijuana, have served in the armed forces, or are a tobacco smoker. These are questions that many dating sites like POF or OKCupid ask for this kind of personal information. Bumble even asks for political affiliation information. This makes sense, users want to find people they are compatible with and this requires personal information. However, users may not trust the company itself to keep the data secure. One might be concerned about [leaking data] to hackers, (https://en.wikipedia.org/wiki/Ashley_Madison_data_breach), [3rd party data sales](http://www.businessinsider.com/spotify-pandora-tinder-apps-sell-anonymized-data-2017-5), or even employees [abusing data access](https://www.theverge.com/2018/1/25/16934064/lyft-customer-data-abuse-allegations). This is data that could effect insurance personal relationships, even employment opportunities.

Using Differential Privacy, we can maintain privacy while allowing for aggregate statistics. Randomized Response permutes a certain percentage of the responses. Since we don't know which were altered, we can't be 100% sure that any response is actually accurate. This gives any user feasible deniability of any data our database has. We can use basic proability theory to reconstruct accruate population statistics from our altered sample. For example the following query will let us know the percent of users who are smokers

```sql
-- alpha = .5 and beta=.5 satisfies a Differential Privacy equal to ln(3)
SELECT
  probabilityrandomresponse(altered_smokers / n, .5, .5) as true_pct_smokers
FROM (
  SELECT sum(randomresponse(is_smoker, .5, .5)) as altered_smokers, count(*) as n
  FROM users
)
```

This works and you can multiply the real probability by `n` to get the real count. However technically the real answers are stored in the database. That still makes this data susceptible to leaks and malicious employees. To really make this private we should only store the altered rendition. However, this would mean potentially showing smokers to people who wanted to filter out smokers. Honestly, I think the false positives, in the app experience, are a worthwile trade for user privacy. 

### Todo

- [ ] Tests
- [ ] Fact Table Examples
- [ ] Laplace Methods for continous variables


