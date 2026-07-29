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
