-- Init user with password admin
INSERT INTO public.users (full_name,email,phone_number,"password",status) VALUES
	 ('Nathanael Cruickshank','Price.Price85@hotmail.com','696-408-1265','$2a$10$cGv95S.sNrdukD3QQmwTkelgz/uW5koN/tihu93Zqz9kcf8SEHAT.',true),
	 ('Lorenza Graham','Katherine_Trantow@hotmail.com','560-671-0415','$2a$10$4WCsekSvfHRLM5/IYuDYjuXtP0RpW5oSQu5.3uDL3/j/b6z978JW6',true),
	 ('Mallie Grant','Shaniya.Mayer33@hotmail.com','674-425-3760','$2a$10$lOD52ZfLjIkMhi37t/Rppe7qWa1t7KcI2Nqy/4FNKkufRMeZ.axnS',true),
	 ('Arielle Gulgowski','Jayda52@yahoo.com','417-927-6456','$2a$10$wUqC/obmk17AMzxeLmychOfA1pjkhDNERzkl87BycDcOeuQhFbT82',true),
	 ('Arvid Hudson','Oliver_Heller81@gmail.com','796-641-9993','$2a$10$0no/8pCvBzP2mR4UM3HdWOaJGrNcObDJghaibt7YgTHpqdMzTTIma',true),
	 ('Erna Ortiz','Drew.Lueilwitz@gmail.com','744-322-0964','$2a$10$mUwSleaEimRTE2Idfuj1l.v3uI14cXHq7jcRNiUBlawPC3/hurRT6',true);

-- Init Crypto symbol
INSERT INTO public.crypto (symbol,status) VALUES
	 ('IDR',true),
	 ('BTC',true),
	 ('ETH',true),
	 ('BNB',true),
	 ('ADA',true),
	 ('DOGE',true),
	 ('XRP',true),
	 ('DOT',true),
	 ('LTC',true),
	 ('LINK',true),
	 ('BCH',true);

-- Init user wallet
INSERT INTO public.wallet (user_id,crypto_id,quantity) VALUES
	 (1,1,3000000.0),
	 (1,2,20.0),
	 (1,5,20000.0),
	 (1,6,30000.0),
	 (1,2,200.0),
	 (2,1,4000000.0),
	 (2,2,25.0),
	 (2,5,20000.0),
	 (2,6,30000.0),
	 (2,2,200.0),
	 (3,1,6000000.0),
	 (3,2,30.0),
	 (3,5,20000.0);
INSERT INTO public.wallet (user_id,crypto_id,quantity) VALUES
	 (3,6,30000.0),
	 (3,2,200.0),
	 (4,1,7000000.0),
	 (4,2,40.0),
	 (4,5,20000.0),
	 (4,6,30000.0),
	 (4,2,200.0),
	 (5,1,12000000.0),
	 (5,2,35.0),
	 (5,5,20000.0),
	 (5,6,30000.0),
	 (5,2,200.0);
INSERT INTO public.wallet (user_id,crypto_id,quantity) VALUES
	 (6,1,9000000.0),
	 (6,2,25.0),
	 (6,5,20000.0),
	 (6,6,30000.0),
	 (6,2,200.0);
