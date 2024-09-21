ALTER table districts
ADD CONSTRAINT fk_province_districts
FOREIGN KEY (province_id) REFERENCES provinces(id)
ON DELETE CASCADE;

ALTER table schools
ADD CONSTRAINT fk_district_schools
FOREIGN KEY (district_id) REFERENCES districts(id)
ON DELETE CASCADE;