import 'i18next';

declare module 'i18next' {
  interface CustomTypeOptions {
    defaultNS: 'translation';
    resources: {
      translation: {
        siteManagement: string;
        addSite: string;
        editSite: string;
        name: string;
        domain: string;
        description: string;
        createdAt: string;
        actions: string;
        edit: string;
        delete: string;
        cancel: string;
        submit: string;
        submitting: string;
        confirmDeleteSite: string;
        siteCreatedSuccess: string;
        siteUpdatedSuccess: string;
        siteDeletedSuccess: string;
        siteCreateFailed: string;
        siteUpdateFailed: string;
        siteDeleteFailed: string;
        loadSitesFailed: string;
        // 添加其他翻译键...
      };
    };
  }
}
