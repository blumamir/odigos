'use client';

import { useQuery } from '@apollo/client/react';
import { GET_GO_ENTERPRISE_OFFSETS } from '@/graphql';

export type GoEnterpriseOffsetMinorVersion = {
  minorVersion: string;
  versions: string[];
};

export type GoEnterpriseOffsetModule = {
  module: string;
  minVersion: string;
  maxVersion: string;
  minorVersions: GoEnterpriseOffsetMinorVersion[];
};

export type GoEnterpriseOffsets = {
  timestamp: string;
  mods: GoEnterpriseOffsetModule[];
};

export type GoEnterpriseOffsetVersionUpdate = {
  version: string;
  isNew: boolean;
};

export type GoEnterpriseOffsetMinorVersionUpdate = {
  minorVersion: string;
  isNew: boolean;
  versions: GoEnterpriseOffsetVersionUpdate[];
};

export type GoEnterpriseOffsetModuleUpdate = {
  module: string;
  isNew: boolean;
  minVersion: string;
  maxVersion: string;
  minorVersions: GoEnterpriseOffsetMinorVersionUpdate[];
};

export type GoEnterpriseOffsetsUpdateCheck = {
  hasUpdates: boolean;
  currentTimestamp: string;
  proposedTimestamp: string;
  mods: GoEnterpriseOffsetModuleUpdate[];
};

type GoEnterpriseOffsetsResponse = {
  goEnterpriseOffsets: GoEnterpriseOffsets;
};

export const useGoEnterpriseOffsets = () => {
  const { data, loading, error, refetch } = useQuery<GoEnterpriseOffsetsResponse>(GET_GO_ENTERPRISE_OFFSETS);

  return {
    goEnterpriseOffsets: data?.goEnterpriseOffsets,
    loading,
    error,
    refetch,
  };
};
