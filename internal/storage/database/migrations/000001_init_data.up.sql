CREATE TABLE gauge_metrics (
    id      VARCHAR PRIMARY KEY,
    value   DOUBLE PRECISION NOT NULL
);

CREATE TABLE counter_metrics (
    id      VARCHAR PRIMARY KEY,
    value   INTEGER NOT NULL
);
