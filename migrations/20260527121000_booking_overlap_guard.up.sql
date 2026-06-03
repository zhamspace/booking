CREATE OR REPLACE FUNCTION booking_prevent_active_overlap()
RETURNS trigger AS $$
BEGIN
    IF NEW.status <> 'cancelled' THEN
        PERFORM pg_advisory_xact_lock(hashtextextended(NEW.venue_id || ':' || NEW.resource_id, 0));

        IF EXISTS (
            SELECT 1
            FROM booking b
            WHERE b.id <> NEW.id
              AND b.venue_id = NEW.venue_id
              AND b.resource_id = NEW.resource_id
              AND b.status <> 'cancelled'
              AND NOT (
                  b.status = 'created'
                  AND b.hold_expires_at IS NOT NULL
                  AND b.hold_expires_at <= now()
              )
              AND tstzrange(b.start_at, b.end_at, '[)') && tstzrange(NEW.start_at, NEW.end_at, '[)')
        ) THEN
            RAISE EXCEPTION 'booking time slot is unavailable'
                USING ERRCODE = '23P01',
                      CONSTRAINT = 'booking_no_active_resource_overlap';
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_booking_prevent_active_overlap ON booking;

CREATE TRIGGER trg_booking_prevent_active_overlap
    BEFORE INSERT OR UPDATE OF venue_id, resource_id, start_at, end_at, status, hold_expires_at
    ON booking
    FOR EACH ROW
    EXECUTE FUNCTION booking_prevent_active_overlap();
