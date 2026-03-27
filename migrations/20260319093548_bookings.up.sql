create table if not exists booking
(
    id                text        default gen_random_uuid() primary key,
    venue_id          text                                   not null,
    resource_id       text                                   not null,
    session_id        text                                   null,
    user_id           text                                   not null,
    status            text        default 'created'          not null,
    payment_status    text        default 'pending'          not null,
    price_total       bigint      default 0                  not null,
    currency          text        default 'KZT'              not null,
    start_at          timestamptz                            not null,
    end_at            timestamptz                            not null,
    timezone          text                                   not null,
    created_at        timestamptz default now()              not null,
    updated_at        timestamptz default now()              not null,
    hold_expires_at   timestamptz                            null,
    payment_intent_id text                                   null,
    cancel_reason     text                                   null,
    cancelled_at      timestamptz                            null,
    confirmed_at      timestamptz                            null,
    completed_at      timestamptz                            null,

    constraint booking_time_range_chk check (end_at > start_at)
);

create index if not exists idx_booking_user_created_at_desc on booking (user_id, created_at desc);
create index if not exists idx_booking_venue_resource_start_at on booking (venue_id, resource_id, start_at);
create index if not exists idx_booking_session_created_at_desc on booking (session_id, created_at desc);
create index if not exists idx_booking_status_created_at_desc on booking (status, created_at desc);
