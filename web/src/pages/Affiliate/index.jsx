/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import React, { useContext, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { API, copy, showError, showSuccess } from '../../helpers';
import { UserContext } from '../../context/User';
import AffiliatePromotionBlock from '../../components/topup/AffiliatePromotionBlock';

const Affiliate = () => {
  const { t } = useTranslation();
  const [, userDispatch] = useContext(UserContext);
  const [affLink, setAffLink] = useState('');
  const [dashboard, setDashboard] = useState(null);

  const getUserQuota = async () => {
    const res = await API.get('/api/user/self');
    const { success, message, data } = res.data;
    if (success) {
      userDispatch({ type: 'login', payload: data });
    } else {
      showError(message);
    }
  };

  const getAffiliateDashboard = async () => {
    const res = await API.get('/api/affiliate/dashboard');
    const { success, message, data } = res.data;
    if (success) {
      setDashboard(data);
      setAffLink(`${window.location.origin}/register?aff=${data.aff_code}`);
    } else {
      showError(message);
    }
  };

  const handleAffLinkClick = async () => {
    await copy(affLink);
    showSuccess(t('邀请链接已复制到剪切板'));
  };

  const handleAffCodeClick = async () => {
    const affCode = affLink ? affLink.split('aff=').pop() : '';
    if (!affCode) {
      showError(t('邀请码为空'));
      return;
    }
    await copy(affCode);
    showSuccess(t('邀请码已复制到剪切板'));
  };

  useEffect(() => {
    getUserQuota().then();
    getAffiliateDashboard().then();
  }, []);

  return (
    <div className='w-full max-w-7xl mx-auto relative min-h-screen lg:min-h-0 mt-[60px] px-2'>
      <AffiliatePromotionBlock
        t={t}
        dashboard={dashboard}
        affLink={affLink}
        handleAffLinkClick={handleAffLinkClick}
        handleAffCodeClick={handleAffCodeClick}
      />
    </div>
  );
};

export default Affiliate;
