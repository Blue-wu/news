// frontend/src/stores/sitemap.js
import { sitemapAPI } from '@/utils/api';

export default {
  namespaced: true,
  state: {
    settings: {
      cacheEnabled: true,
      cacheDuration: 86400,
      cacheStatus: {
        isCached: false,
        cacheAge: '',
        cachedAt: ''
      }
    },
    statistics: {
      visitCount: 0,
      lastVisitTime: ''
    },
    isLoading: false,
    error: null
  },
  
  mutations: {
    SET_SETTINGS(state, settings) {
      state.settings = settings;
    },
    
    SET_STATISTICS(state, statistics) {
      state.statistics = statistics;
    },
    
    SET_LOADING(state, status) {
      state.isLoading = status;
    },
    
    SET_ERROR(state, error) {
      state.error = error;
    }
  },
  
  actions: {
    async fetchSettings({ commit }) {
      commit('SET_LOADING', true);
      try {
        const response = await sitemapAPI.getSettings();
        if (response.data.success) {
          commit('SET_SETTINGS', response.data.settings);
          commit('SET_STATISTICS', response.data.statistics);
        }
        return response.data;
      } catch (error) {
        commit('SET_ERROR', error.message);
        throw error;
      } finally {
        commit('SET_LOADING', false);
      }
    },
    
    async updateSettings({ commit }, settings) {
      commit('SET_LOADING', true);
      try {
        const response = await sitemapAPI.updateSettings(settings);
        if (response.data.success) {
          // 重新获取设置以更新状态
          const freshResponse = await sitemapAPI.getSettings();
          if (freshResponse.data.success) {
            commit('SET_SETTINGS', freshResponse.data.settings);
          }
        }
        return response.data;
      } catch (error) {
        commit('SET_ERROR', error.message);
        throw error;
      } finally {
        commit('SET_LOADING', false);
      }
    },
    
    async regenerateSitemap({ commit }) {
      commit('SET_LOADING', true);
      try {
        const response = await sitemapAPI.regenerate();
        if (response.data.success) {
          // 重新获取设置以更新状态
          const freshResponse = await sitemapAPI.getSettings();
          if (freshResponse.data.success) {
            commit('SET_SETTINGS', freshResponse.data.settings);
          }
        }
        return response.data;
      } catch (error) {
        commit('SET_ERROR', error.message);
        throw error;
      } finally {
        commit('SET_LOADING', false);
      }
    },
    
    async clearCache({ commit }) {
      commit('SET_LOADING', true);
      try {
        const response = await sitemapAPI.clearCache();
        if (response.data.success) {
          // 重新获取设置以更新状态
          const freshResponse = await sitemapAPI.getSettings();
          if (freshResponse.data.success) {
            commit('SET_SETTINGS', freshResponse.data.settings);
          }
        }
        return response.data;
      } catch (error) {
        commit('SET_ERROR', error.message);
        throw error;
      } finally {
        commit('SET_LOADING', false);
      }
    }
  }
};