#!/usr/bin/env python3
"""
Fix :data binding mismatches in views/members Vue files.
The template uses :data="data" but script uses different variable names.
This script normalizes all to use the actual data variable name from the script.
"""
import os
import re
import glob

MEMBERS_DIR = r"C:\Users\YKing\.openclaw\workspace\mdm-project\frontend\src\views\members"

# Map of files that use non-standard data variable names
# key: filename, value: the actual data variable name used in the script
DATA_VAR_MAP = {
    "AmountDiscountView.vue": "dataList",
    "AmountReduceView.vue": "dataList",
    "BuyGiftView.vue": "dataList",
    "CouponGrantView.vue": "dataList",
    "CouponInventoryView.vue": "dataList",
    "CouponMessagesView.vue": "dataList",
    "DirectReduceView.vue": "dataList",
    "GiftRecordsView.vue": "dataList",
    "HighFreqTagView.vue": "dataList",
    "InterestTagView.vue": "dataList",
    "LowFreqTagView.vue": "dataList",
    "MemberArticlesView.vue": "dataList",
    "MemberCardGroupsView.vue": "dataList",
    "MemberCardTypesView.vue": "dataList",
    "MemberCardView.vue": "cardTypeList",
    "MemberCoupons.vue": "couponList",
    "MemberLevelRulesView.vue": "dataList",
    "MemberOrdersView.vue": "dataList",
    "OccupationTypesView.vue": "dataList",
    "PointsExcludeView.vue": "dataList",
    "PointsInventory.vue": "data",
    "PointsRecords.vue": "records",
    "PointsRules.vue": "dataList",
    "PromotionTypesView.vue": "dataList",
    "RedpacketView.vue": "dataList",
    "SmsChannelView.vue": "dataList",
    "SmsTemplateView.vue": "dataList",
    "TagView.vue": "tagList",
    "VipExclusiveView.vue": "dataList",
    # Files that also need their script data var normalized to "data"
    "PointsExclude.vue": None,  # Need to rename rules to data
    "MemberBenefitsView.vue": None,  # Need to rename data to data
    "MemberCardGroupsView.vue": None,  # Already in map
    "MiniProgram.vue": None,  # Need to check
    "MemberStores.vue": None,  # Need to check
    "StoreView.vue": None,  # Need to check
    "MemberDetailView.vue": None,  # Need to check
    "MemberPoints.vue": None,  # Need to check
    "MemberPromotions.vue": None,  # Need to check
    "MemberReceptionView.vue": None,  # Need to check
    "MemberSettingsView.vue": None,  # Need to check
    "PointsView.vue": None,  # Need to check
    "TagAutoCleanView.vue": None,  # Need to check
    "TagReportView.vue": None,  # Need to check
}

def extract_script_data_var(content):
    """Extract the primary data ref variable name from script section."""
    # Find script section
    script_match = re.search(r'<script[^>]*>(.*?)</script>', content, re.DOTALL)
    if not script_match:
        return None
    
    script = script_match.group(1)
    
    # Look for data ref patterns
    patterns = [
        r'const\s+(data)\s*=\s*ref\(',           # const data = ref(
        r'const\s+(dataList)\s*=\s*ref\(',         # const dataList = ref(
        r'const\s+(records)\s*=\s*ref\(',           # const records = ref(
        r'const\s+(rules)\s*=\s*ref\(',             # const rules = ref(
        r'const\s+(couponList)\s*=\s*ref\(',        # const couponList = ref(
        r'const\s+(memberList)\s*=\s*ref\(',        # const memberList = ref(
        r'const\s+(levelList)\s*=\s*ref\(',         # const levelList = ref(
        r'const\s+(tagList)\s*=\s*ref\(',           # const tagList = ref(
        r'const\s+(stores)\s*=\s*ref\(',             # const stores = ref(
        r'const\s+(member)\s*=\s*ref\(',           # const member = ref(
        r'const\s+(cardTypeList)\s*=\s*ref\(',     # const cardTypeList = ref(
        r'const\s+(activeTab)\s*=\s*ref\(',        # const activeTab = ref(
    ]
    
    for pattern in patterns:
        m = re.search(pattern, script)
        if m:
            return m.group(1)
    
    return None

def fix_file(filepath):
    """Fix a single file's data binding."""
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    filename = os.path.basename(filepath)
    
    # Extract the actual data variable from the script
    actual_data_var = extract_script_data_var(content)
    
    if not actual_data_var:
        print(f"  [SKIP] {filename}: no data ref found")
        return False
    
    if actual_data_var == 'data':
        print(f"  [OK]   {filename}: already uses 'data'")
        return True
    
    # Fix the template: replace :data="data" with :data="<actual_var>"
    old_binding = ':data="data"'
    new_binding = f':data="{actual_data_var}"'
    
    if old_binding in content:
        content = content.replace(old_binding, new_binding)
        print(f"  [FIX]  {filename}: changed :data=\"data\" -> :data=\"{actual_data_var}\"")
    else:
        print(f"  [OK]   {filename}: no :data=\"data\" binding found (uses :data=\"{actual_data_var}\"?)")
    
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)
    
    return True

def main():
    files = glob.glob(os.path.join(MEMBERS_DIR, "*.vue"))
    
    # Filter out non-Vue files and the script itself
    files = [f for f in files if not f.endswith('rewrite_members.py') and not f.endswith('rewrite_batch.py') and not f.endswith('fixall.cjs')]
    
    print(f"Processing {len(files)} Vue files...")
    
    results = {'fixed': [], 'skipped': [], 'errors': []}
    
    for filepath in sorted(files):
        filename = os.path.basename(filepath)
        try:
            fixed = fix_file(filepath)
            if fixed:
                results['fixed'].append(filename)
            else:
                results['skipped'].append(filename)
        except Exception as e:
            results['errors'].append((filename, str(e)))
            print(f"  [ERR]  {filename}: {e}")
    
    print(f"\nSummary:")
    print(f"  Fixed/Skipped: {len(results['fixed'])}")
    print(f"  Errors: {len(results['errors'])}")
    for fname, err in results['errors']:
        print(f"    {fname}: {err}")

if __name__ == '__main__':
    main()
