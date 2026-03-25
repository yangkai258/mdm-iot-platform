-- disease_patterns 种子数据
INSERT INTO disease_patterns (pattern_id, pattern_name, disease_name, species, symptoms, severity, description) VALUES
('dp_001', '猫感冒', 'feline_upper_respiratory', 'cat', '["sneezing", "runny_nose", "lethargy"]'::jsonb, 'mild', '猫咪上呼吸道感染'),
('dp_002', '狗皮肤病', 'canine_dermatitis', 'dog', '["itching", "hair_loss", "redness"]'::jsonb, 'moderate', '狗狗皮肤炎症');
