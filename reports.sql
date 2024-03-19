select
    count(*) as total,
    count(*) filter (where web_app_access_token is not null) as mined,
        count(*) filter (where balance = 5000) as with_start_balance,
        count(*) filter (where stopped_at is not null) as stopped,
        count(*) filter (where stopped_at is null) as active
from users join promos on users.promo_id = promos.id
where promo_id = 3;